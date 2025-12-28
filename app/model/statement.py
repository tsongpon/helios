from typing import List, Literal, Optional

from pydantic import BaseModel, Field


class Transaction(BaseModel):
    """Represents a single credit card transaction."""

    date: str = Field(description="Transaction date in YYYY-MM-DD format")
    description: str = Field(description="Transaction description")
    amount: float = Field(description="Transaction amount (positive number)")
    transaction_type: Literal["debit", "credit"] = Field(
        description="Type of transaction: 'debit' for purchases/charges or 'credit' for payments/refunds"
    )
    is_installment: bool = Field(
        default=False, description="True if this is an installment transaction"
    )
    installment_current: Optional[int] = Field(
        default=None, description="Current installment number (e.g., 3 in '3/12')"
    )
    installment_total: Optional[int] = Field(
        default=None, description="Total number of installments (e.g., 12 in '3/12')"
    )
    installment_plan_amount: Optional[float] = Field(
        default=None, description="Total amount of the installment plan"
    )

    def to_dict(self):
        """Convert to dictionary for backward compatibility."""
        return self.model_dump()


class StatementSummary(BaseModel):
    """Represents the summary section of a credit card statement."""

    card_number: Optional[str] = Field(
        default=None, description="Card number (last 4 digits or masked format)"
    )
    credit_limit: Optional[float] = Field(default=None, description="Credit limit")
    outstanding_balance: float = Field(description="Outstanding balance")
    minimum_payment: Optional[float] = Field(
        default=None, description="Minimum payment amount"
    )
    payment_due_date: Optional[str] = Field(
        default=None, description="Payment due date in YYYY-MM-DD format"
    )
    statement_date: str = Field(description="Statement date in YYYY-MM-DD format")
    previous_balance: Optional[float] = Field(
        default=None, description="Previous balance"
    )
    total_credits: Optional[float] = Field(
        default=None, description="Total credits for the period"
    )
    total_debits: Optional[float] = Field(
        default=None, description="Total debits for the period"
    )

    def to_dict(self):
        """Convert to dictionary for backward compatibility."""
        return self.model_dump()


class CreditCardStatement(BaseModel):
    """Represents a complete credit card statement with summary and transactions."""

    summary: StatementSummary = Field(description="Statement summary information")
    transactions: List[Transaction] = Field(
        description="List of all transactions in the statement"
    )

    def to_dict(self):
        """Convert to dictionary."""
        return {
            "summary": self.summary.to_dict(),
            "transactions": [t.to_dict() for t in self.transactions],
        }

    def to_json(self):
        """Convert to JSON string."""
        return self.model_dump_json(indent=2)
