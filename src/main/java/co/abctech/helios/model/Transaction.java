package co.abctech.helios.model;

import java.time.LocalDate;

public record Transaction(
    LocalDate transactionDate,
    LocalDate postingDate,
    String description,
    Double amount,
    Boolean isInstallment,
    Integer installmentCurrent,
    Integer installmentTotal,
    Double installmentPlanAmount
) {}
