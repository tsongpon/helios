from google import genai
import json
import os
from dotenv import load_dotenv
from model.statement import Transaction, StatementSummary, CreditCardStatement

# Load environment variables from .env file
load_dotenv()


class LlmRepository(object):
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

    def retrieve_statement(self, statementText):
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
            "transaction_type": <"debit" or "credit">
        }}
    ]
}}

Instructions:
1. Extract all transactions with their dates, descriptions, and amounts
2. Identify whether each transaction is a debit (purchase/charge) or credit (payment/refund)
3. Extract summary information like card number (last 4 digits or masked format), credit limit, outstanding balance, payment due date, etc.
4. Convert all dates to YYYY-MM-DD format
5. Ensure all amounts are positive numbers
6. If a field is not found in the statement, use null
7. Return ONLY the JSON object, no additional text or explanation
"""

        try:
            response = self.client.models.generate_content(
                model="gemini-2.0-flash-exp", contents=prompt
            )
            response_text = response.text.strip()

            # Remove markdown code blocks if present
            if response_text.startswith("```json"):
                response_text = response_text[7:]
            if response_text.startswith("```"):
                response_text = response_text[3:]
            if response_text.endswith("```"):
                response_text = response_text[:-3]
            response_text = response_text.strip()

            # Parse JSON response
            data = json.loads(response_text)

            # Create structured data model
            summary = StatementSummary(**data["summary"])
            transactions = [Transaction(**t) for t in data["transactions"]]
            statement = CreditCardStatement(summary=summary, transactions=transactions)

            return statement

        except json.JSONDecodeError as e:
            raise Exception(
                f"Failed to parse LLM response as JSON: {str(e)}\nResponse: {response_text}"
            )
        except Exception as e:
            raise Exception(f"Error processing statement with LLM: {str(e)}")
