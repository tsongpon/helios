from google import genai
import os
import json
import pdfplumber
from dotenv import load_dotenv
from model.statement import Transaction, StatementSummary, CreditCardStatement

# Load environment variables from .env file
load_dotenv()


class CcStatementRepository(object):
    def __init__(self, api_key=None):
        """
        Initialize LLM Repository with Gemini API.

        Args:
            api_key (str): Google API key. If not provided, will look for GOOGLE_API_KEY env variable
        """
        self.api_key = api_key or os.getenv("GOOGLE_API_KEY")
        if not self.api_key:
            raise ValueError(
                "Google API key is required. Set GOOGLE_API_KEY environment variable or pass api_key to constructor."
            )
        self.client = genai.Client(api_key=self.api_key)

    def retrieve_statement(self, file_path, password=None):
        """
        Parse a PDF file and extract all text content.

        Args:
            file_path (str): Path to the PDF file
            password (str, optional): Password for encrypted PDF files

        Returns:
            str: Extracted text from the PDF
        """
        try:
            text = ""

            with pdfplumber.open(file_path, password=password) as pdf:
                for page in pdf.pages:
                    page_text = page.extract_text()
                    if page_text:
                        text += page_text

            return self.parse_by_llm(text)
        except Exception as e:
            raise Exception(f"Error parsing PDF: {str(e)}")

    def parse_by_llm(self, statementText):
        """
        Process credit card statement text using Gemini LLM and extract structured data.

        Args:
            statementText (str): Raw text extracted from credit card statement PDF

        Returns:
            CreditCardStatement: Structured data model containing summary and transactions
        """
        prompt = f"""
Analyze the following credit card statement text and extract structured information.

Statement Text:
{statementText}

Please extract and return the following information in JSON format:
{{
    "summary": {{
        "card_number": <string or null>,
        "credit_limit": <float or null>,
        "outstanding_balance": <float>,
        "minimum_payment": <float or null>,
        "payment_due_date": <string in YYYY-MM-DD format or null>,
        "statement_date": <string in YYYY-MM-DD format>,
        "previous_balance": <float or null>,
        "total_credits": <float or null>,
        "total_debits": <float or null>
    }},
    "transactions": [
        {{
            "date": <string in YYYY-MM-DD format>,
            "description": <string>,
            "amount": <float>,
            "transaction_type": <"debit" or "credit">,
            "is_installment": <boolean>,
            "installment_current": <integer or null>,
            "installment_total": <integer or null>,
            "installment_plan_amount": <float or null>
        }}
    ]
}}

Instructions:
1. Extract all transactions with their dates, descriptions, and amounts
2. Identify whether each transaction is a debit (purchase/charge) or credit (payment/refund)
3. Detect installment transactions - look for patterns like "3/12", "installment 5 of 10", "EMI", "IPP", etc.
4. For installment transactions:
   - Set is_installment to true
   - Extract installment_current (the current payment number, e.g., 3 from "3/12")
   - Extract installment_total (total number of payments, e.g., 12 from "3/12")
   - If available, extract installment_plan_amount (the total purchase amount being financed)
5. For non-installment transactions, set is_installment to false and installment fields to null
6. Extract summary information like card number (last 4 digits or masked format), credit limit, outstanding balance, payment due date, etc.
7. Convert all dates to YYYY-MM-DD format
8. IMPORTANT - Year handling for transaction dates:
   - If transaction dates in the PDF include the year, use that year
   - If transaction dates are missing the year (e.g., only "01/15" or "Jan 15"), infer the year from the statement_date
   - For transactions that occur in the statement period, use the same year as the statement date
   - For transactions near year-end: if the transaction month is December and statement month is January, the transaction year should be (statement year - 1)
   - For transactions near year-start: if the transaction month is January and statement month is December, the transaction year should be (statement year + 1)
9. Ensure all amounts are positive numbers
10. If a field is not found in the statement, use null
11. Return ONLY the JSON object, no additional text or explanation
"""

        try:
            response = self.client.models.generate_content(
                model="gemini-2.0-flash-exp", contents=prompt
            )
            response_text = response.text.strip()

            # Clean up the response text to extract JSON
            cleaned_text = self._extract_json_from_response(response_text)

            # Parse JSON response
            data = json.loads(cleaned_text)

            # Validate required fields
            if "summary" not in data or "transactions" not in data:
                raise ValueError(
                    "Response missing required fields: 'summary' or 'transactions'"
                )

            # Create structured data model
            summary = StatementSummary(**data["summary"])
            transactions = [Transaction(**t) for t in data["transactions"]]
            statement = CreditCardStatement(summary=summary, transactions=transactions)

            return statement

        except json.JSONDecodeError as e:
            # Log the problematic response for debugging
            error_msg = f"Failed to parse LLM response as JSON: {str(e)}\n"
            error_msg += f"Attempted to parse: {cleaned_text[:500]}..."
            raise Exception(error_msg)
        except KeyError as e:
            raise Exception(f"Missing required field in LLM response: {str(e)}")
        except TypeError as e:
            raise Exception(f"Invalid data type in LLM response: {str(e)}")
        except Exception as e:
            raise Exception(f"Error processing statement with LLM: {str(e)}")

    def _extract_json_from_response(self, response_text):
        """
        Extract JSON from LLM response, handling various formatting issues.

        Args:
            response_text (str): Raw response from LLM

        Returns:
            str: Cleaned JSON string
        """
        text = response_text.strip()

        # Remove markdown code blocks
        if text.startswith("```json"):
            text = text[7:]
        elif text.startswith("```"):
            text = text[3:]

        if text.endswith("```"):
            text = text[:-3]

        text = text.strip()

        # Try to find JSON object boundaries if text contains extra content
        if not text.startswith("{"):
            # Look for first occurrence of {
            start_idx = text.find("{")
            if start_idx != -1:
                text = text[start_idx:]

        if not text.endswith("}"):
            # Look for last occurrence of }
            end_idx = text.rfind("}")
            if end_idx != -1:
                text = text[: end_idx + 1]

        return text
