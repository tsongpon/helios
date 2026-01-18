package co.abctech.helios.repository;

import co.abctech.helios.model.Statement;

public interface LlmRepository {
    Statement parseStatement(String statementText);
}
