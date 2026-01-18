package co.abctech.helios.service;

import co.abctech.helios.model.Statement;
import co.abctech.helios.repository.LlmRepository;
import org.springframework.stereotype.Service;

@Service
public class StatementService {

    private final LlmRepository llmRepository;

    public StatementService(LlmRepository llmRepository) {
        this.llmRepository = llmRepository;
    }

    public Statement parseStatement(String statementText) {
        if (statementText == null || statementText.trim().isEmpty()) {
            throw new IllegalArgumentException(
                "Statement text cannot be null or empty"
            );
        }

        return llmRepository.parseStatement(statementText);
    }
}
