import json
from dataclasses import dataclass, asdict
from typing import List, Optional


@dataclass
class Transaction:
    date: str
    description: str
    amount: float
    transaction_type: str  # "debit" or "credit"
    is_installment: bool = False  # True if this is an installment transaction
    installment_current: Optional[int] = (
        None  # Current installment number (e.g., 3 in "3/12")
    )
    installment_total: Optional[int] = (
        None  # Total number of installments (e.g., 12 in "3/12")
    )
    installment_plan_amount: Optional[float] = (
        None  # Total amount of the installment plan
    )

    def to_dict(self):
        return asdict(self)


@dataclass
class StatementSummary:
    card_number: Optional[str]
    credit_limit: Optional[float]
    outstanding_balance: float
    minimum_payment: Optional[float]
    payment_due_date: Optional[str]
    statement_date: str
    previous_balance: Optional[float]
    total_credits: Optional[float]
    total_debits: Optional[float]

    def to_dict(self):
        return asdict(self)


@dataclass
class CreditCardStatement:
    summary: StatementSummary
    transactions: List[Transaction]

    def to_dict(self):
        return {
            "summary": self.summary.to_dict(),
            "transactions": [t.to_dict() for t in self.transactions],
        }

    def to_json(self):
        return json.dumps(self.to_dict(), ensure_ascii=False, indent=2)
