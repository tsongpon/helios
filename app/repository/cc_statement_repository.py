import os

import pdfplumber
from dotenv import load_dotenv
from langchain_core.output_parsers import PydanticOutputParser
from langchain_core.prompts import ChatPromptTemplate
from langchain_google_genai import ChatGoogleGenerativeAI

from model.statement import CreditCardStatement

# Load environment variables from .env file
load_dotenv()


class CcStatementRepository(object):
    def __init__(self, api_key=None):
        """
        Initialize Credit Card Statement Repository with Gemini API using LangChain.

        Args:
            api_key (str): Google API key. If not provided, will look for GOOGLE_API_KEY env variable
        """
        self.api_key = api_key or os.getenv("GOOGLE_API_KEY")
        if not self.api_key:
            raise ValueError(
                "Google API key is required. Set GOOGLE_API_KEY environment variable or pass api_key to constructor."
            )

        # Initialize LangChain LLM
        self.llm = ChatGoogleGenerativeAI(
            model="gemini-2.0-flash-exp",
            api_key=self.api_key,
            temperature=0,  # Use deterministic output for structured data extraction
        )

        # Initialize output parser
        self.parser = PydanticOutputParser(pydantic_object=CreditCardStatement)

    def retrieve_statement(self, file_path, password=None):
        """
        Parse a PDF file and extract structured credit card statement data.

        Args:
            file_path (str): Path to the PDF file
            password (str, optional): Password for encrypted PDF files

        Returns:
            CreditCardStatement: Structured statement data
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
        # Create prompt template
        prompt = ChatPromptTemplate.from_template(
            """Analyze the following credit card statement text and extract structured information.

Statement Text:
{statement_text}

{format_instructions}

CRITICAL INSTRUCTIONS - MUST FOLLOW EXACTLY:

1. Extract all transactions with their dates, descriptions, and amounts

2. Identify whether each transaction is a debit (purchase/charge) or credit (payment/refund)

3. **INSTALLMENT DETECTION - VERY IMPORTANT**:
   - Carefully examine EACH transaction description for installment patterns
   - Common installment patterns to look for:
     * "XX/YY" format (e.g., "09/10", "03/06", "10/10") - this is the MOST COMMON pattern
     * "X of Y" format (e.g., "3 of 12", "installment 5 of 10")
     * Keywords: "IPP", "EMI", "INST", "INSTALLMENT"
   - The installment pattern is typically at the END of the description
   - Example: "ZOOM CAMERA-WEST GATE 09/10" means installment 9 of 10
   - Example: "2C2P *LAZADA 03/06" means installment 3 of 6

4. For installment transactions (when you find XX/YY or similar patterns):
   - **MUST** set is_installment to true
   - **MUST** extract installment_current (the first number, e.g., 9 from "09/10")
   - **MUST** extract installment_total (the second number, e.g., 10 from "09/10")
   - If the total plan amount is shown separately, extract installment_plan_amount
   - Otherwise, leave installment_plan_amount as null

5. For non-installment transactions:
   - Set is_installment to false
   - Set installment_current to null
   - Set installment_total to null
   - Set installment_plan_amount to null

6. Extract summary information like card number (last 4 digits or masked format), credit limit, outstanding balance, payment due date, etc.

7. Convert all dates to YYYY-MM-DD format

8. Year handling for transaction dates:
   - If transaction dates in the PDF include the year, use that year
   - If transaction dates are missing the year (e.g., only "01/15" or "Jan 15"), infer the year from the statement_date
   - For transactions that occur in the statement period, use the same year as the statement date
   - For transactions near year-end: if the transaction month is December and statement month is January, the transaction year should be (statement year - 1)
   - For transactions near year-start: if the transaction month is January and statement month is December, the transaction year should be (statement year + 1)

9. Ensure all amounts are positive numbers

10. If a field is not found in the statement, use null

REMINDER: DO NOT MISS ANY INSTALLMENT PATTERNS! Check every transaction description carefully for XX/YY patterns."""
        )

        try:
            # Create the chain: prompt -> LLM -> parser
            chain = prompt | self.llm | self.parser

            statement = chain.invoke(
                {
                    "statement_text": statementText,
                    "format_instructions": self.parser.get_format_instructions(),
                }
            )

            return statement

        except Exception as e:
            raise Exception(f"Error processing statement with LLM: {str(e)}")
