from repository.cc_statement_repository import CcStatementRepository


def main():
    repo = CcStatementRepository()
    resp = repo.retrieve_statement(file_path="statement.pdf", password="15June1981")
    print(resp.model_dump_json())


if __name__ == "__main__":
    main()
