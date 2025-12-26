from repository.cc_statement_repository import CcStatementRepository


def main():
    repo = CcStatementRepository(api_key="AIzaSyAQtt2xrDDmEkPVwxrAadE-w0-yukgUnDI")
    text = repo.parse_pdf(file="statement.pdf", password="15June1981")
    resp = repo.retrieve_statement(text)
    print(resp.to_json())


if __name__ == "__main__":
    main()
