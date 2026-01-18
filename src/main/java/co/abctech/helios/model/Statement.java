package co.abctech.helios.model;

import java.time.LocalDate;
import java.util.List;

public record Statement(
    String cardNumber,
    LocalDate statementDate,
    Double creditLine,
    Double totalPaymentDue,
    Double outstandingBalance,
    Double minimumPaymentDue,
    Double availableCredit,
    LocalDate paymentDueDate,
    List<Transaction> transactions
) {}
