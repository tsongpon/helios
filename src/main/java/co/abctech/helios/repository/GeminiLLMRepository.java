package co.abctech.helios.repository;

import co.abctech.helios.model.Statement;
import co.abctech.helios.model.Transaction;
import dev.langchain4j.model.chat.ChatLanguageModel;
import dev.langchain4j.model.googleai.GoogleAiGeminiChatModel;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.time.format.DateTimeParseException;
import java.util.ArrayList;
import java.util.List;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Repository;

@Repository
public class GeminiLLMRepository implements LlmRepository {

    private final ChatLanguageModel chatModel;

    public GeminiLLMRepository(
        @Value("${gemini.api-key}") String apiKey,
        @Value("${gemini.model-name:gemini-1.5-flash}") String modelName
    ) {
        this.chatModel = GoogleAiGeminiChatModel.builder()
            .apiKey(apiKey)
            .modelName(modelName)
            .temperature(0.0)
            .maxOutputTokens(8000)
            .maxRetries(1)
            .logRequestsAndResponses(true)
            .build();
    }

    @Override
    public Statement parseStatement(String statementText) {
        try {
            String prompt = buildPrompt(statementText);
            String response = chatModel.generate(prompt);
            return parseResponse(response);
        } catch (Exception e) {
            throw new RuntimeException(
                "Failed to parse statement using Gemini: " + e.getMessage(),
                e
            );
        }
    }

    private String buildPrompt(String statementText) {
        return """
        Extract credit card statement data. Return in this EXACT pipe-delimited format:

        HEADER|cardNumber|statementDate|creditLine|totalPaymentDue|outstandingBalance|minimumPaymentDue|availableCredit|paymentDueDate
        TXN|transactionDate|postingDate|description|amount|isInstallment|installmentCurrent|installmentTotal|installmentPlanAmount
        TXN|transactionDate|postingDate|description|amount|isInstallment|installmentCurrent|installmentTotal|installmentPlanAmount
        ...

        Rules:
        - First line starts with HEADER| followed by 8 values separated by |
        - Each transaction line starts with TXN| followed by 8 values separated by |
        - Dates in DD/MM/YY format (infer year from paymentDueDate if only DD/MM shown)
        - Numbers without commas
        - Descriptions max 25 chars
        - isInstallment: true or false
        - Use empty string for null values
        - outstandingBalance: look for "NEW BALANCE"/"BALANCE"/"THIS PERIOD BALANCE"
        - availableCredit: calculate as creditLine - outstandingBalance
        - Include ALL transactions

        Example:
        HEADER|5468 48XX XXXX 8032|05/10/23|600000|51011.75|69386.75|12779.48|530613.25|25/10/23
        TXN|05/09/23|06/09/23|WWW.GRAB.COM|125|false|||
        TXN|09/05/23|05/10/23|SHOPEEPAY INSTALLMENT|22050|true|5|10|3675

        Statement:
        %s
        """.formatted(statementText);
    }

    private Statement parseResponse(String response) {
        try {
            // Clean response - remove markdown code blocks if present
            String cleanedResponse = response
                .trim()
                .replaceAll("^```.*\n", "")
                .replaceAll("```\\s*$", "")
                .trim();

            // Parse pipe-delimited format
            String[] lines = cleanedResponse.split("\n");
            StatementDto header = null;
            List<TransactionDto> transactions = new ArrayList<>();

            for (String line : lines) {
                line = line.trim();
                if (line.isEmpty()) continue;

                if (line.startsWith("HEADER|")) {
                    header = parseHeaderLine(line);
                } else if (line.startsWith("TXN|")) {
                    transactions.add(parseTransactionLine(line));
                }
            }

            if (header == null) {
                throw new RuntimeException(
                    "No HEADER line found in LLM response"
                );
            }

            header.transactions = transactions;
            return convertToStatement(header);
        } catch (Exception e) {
            throw e;
        }
    }

    private StatementDto parseHeaderLine(String line) {
        String[] parts = line.substring(7).split("\\|", -1); // Remove "HEADER|" and split
        StatementDto dto = new StatementDto();
        dto.cardNumber = parts[0].trim().isEmpty() ? null : parts[0].trim();
        dto.statementDate = parts[1].trim().isEmpty() ? null : parts[1].trim();
        dto.creditLine = parseDouble(parts[2]);
        dto.totalPaymentDue = parseDouble(parts[3]);
        dto.outstandingBalance = parseDouble(parts[4]);
        dto.minimumPaymentDue = parseDouble(parts[5]);
        dto.availableCredit = parseDouble(parts[6]);
        dto.paymentDueDate = parts[7].trim().isEmpty() ? null : parts[7].trim();
        return dto;
    }

    private TransactionDto parseTransactionLine(String line) {
        String[] parts = line.substring(4).split("\\|", -1); // Remove "TXN|" and split
        TransactionDto dto = new TransactionDto();
        dto.transactionDate = parts[0].trim().isEmpty()
            ? null
            : parts[0].trim();
        dto.postingDate = parts[1].trim().isEmpty() ? null : parts[1].trim();
        dto.description = parts[2].trim().isEmpty() ? null : parts[2].trim();
        dto.amount = parseDouble(parts[3]);
        dto.isInstallment = parseBoolean(parts[4]);
        dto.installmentCurrent = parseInteger(parts[5]);
        dto.installmentTotal = parseInteger(parts[6]);
        dto.installmentPlanAmount = parseDouble(parts[7]);
        return dto;
    }

    private Double parseDouble(String value) {
        if (value == null || value.trim().isEmpty()) return null;
        try {
            return Double.parseDouble(value.trim());
        } catch (NumberFormatException e) {
            return null;
        }
    }

    private Integer parseInteger(String value) {
        if (value == null || value.trim().isEmpty()) return null;
        try {
            return Integer.parseInt(value.trim());
        } catch (NumberFormatException e) {
            return null;
        }
    }

    private Boolean parseBoolean(String value) {
        if (value == null || value.trim().isEmpty()) return null;
        return Boolean.parseBoolean(value.trim());
    }

    private Statement convertToStatement(StatementDto dto) {
        List<Transaction> transactions = new ArrayList<>();
        if (dto.transactions != null) {
            for (TransactionDto txDto : dto.transactions) {
                transactions.add(
                    new Transaction(
                        parseDate(txDto.transactionDate),
                        parseDate(txDto.postingDate),
                        txDto.description,
                        txDto.amount,
                        txDto.isInstallment,
                        txDto.installmentCurrent,
                        txDto.installmentTotal,
                        txDto.installmentPlanAmount
                    )
                );
            }
        }

        return new Statement(
            dto.cardNumber,
            parseDate(dto.statementDate),
            dto.creditLine,
            dto.totalPaymentDue,
            dto.outstandingBalance,
            dto.minimumPaymentDue,
            dto.availableCredit,
            parseDate(dto.paymentDueDate),
            transactions
        );
    }

    private LocalDate parseDate(String dateStr) {
        if (dateStr == null || dateStr.trim().isEmpty()) {
            return null;
        }

        try {
            // Try DD/MM/YY format (e.g., "05/08/25")
            DateTimeFormatter formatter = DateTimeFormatter.ofPattern(
                "dd/MM/yy"
            );
            return LocalDate.parse(dateStr.trim(), formatter);
        } catch (DateTimeParseException e) {
            try {
                // Try DD/MM/YYYY format as fallback
                DateTimeFormatter formatter = DateTimeFormatter.ofPattern(
                    "dd/MM/yyyy"
                );
                return LocalDate.parse(dateStr.trim(), formatter);
            } catch (DateTimeParseException e2) {
                return null;
            }
        }
    }

    // DTOs for JSON deserialization
    private static class StatementDto {

        String cardNumber;
        String statementDate;
        Double creditLine;
        Double totalPaymentDue;
        Double outstandingBalance;
        Double minimumPaymentDue;
        Double availableCredit;
        String paymentDueDate;
        List<TransactionDto> transactions;
    }

    private static class TransactionDto {

        String transactionDate;
        String postingDate;
        String description;
        Double amount;
        Boolean isInstallment;
        Integer installmentCurrent;
        Integer installmentTotal;
        Double installmentPlanAmount;
    }
}
