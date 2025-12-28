import os
import tempfile

from fastapi import FastAPI, File, Form, HTTPException, UploadFile
from fastapi.responses import JSONResponse
from repository.cc_statement_repository import CcStatementRepository

app = FastAPI(title="Credit Card Statement Parser API")

# Initialize repositories
cc_repo = CcStatementRepository()


@app.post("/api/parse-statement")
async def parse_statement(
    file: UploadFile = File(..., description="PDF file of the credit card statement"),
    password: str = Form(None, description="Password for encrypted PDF (optional)"),
):
    """
    Upload a credit card statement PDF and get structured data.

    Args:
        file: PDF file to upload
        password: Optional password if the PDF is encrypted

    Returns:
        JSON response with statement summary and transactions
    """
    # Validate file type
    if not file.filename.endswith(".pdf"):
        raise HTTPException(status_code=400, detail="Only PDF files are accepted")

    # Create a temporary file to store the uploaded PDF
    temp_file = None
    try:
        # Save uploaded file to temporary location
        with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf") as temp_file:
            content = await file.read()
            temp_file.write(content)
            temp_file_path = temp_file.name

        # Parse PDF to text
        try:
            statement = cc_repo.retrieve_statement(temp_file_path, password=password)
        except Exception as e:
            if "password" in str(e).lower() or "encrypt" in str(e).lower():
                raise HTTPException(
                    status_code=401,
                    detail="PDF is password-protected. Please provide the correct password.",
                )
            raise HTTPException(status_code=400, detail=f"Error parsing PDF: {str(e)}")

        # Return structured data as JSON
        return JSONResponse(content=statement.to_dict(), status_code=200)

    except HTTPException:
        raise
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Unexpected error: {str(e)}")
    finally:
        # Clean up temporary file
        if temp_file and os.path.exists(temp_file_path):
            os.unlink(temp_file_path)


@app.get("/")
async def root():
    """
    Root endpoint with API information.
    """
    return {
        "message": "Credit Card Statement Parser API",
        "version": "1.0.0",
        "endpoints": {
            "/api/parse-statement": "POST - Upload PDF and get structured statement data"
        },
    }


@app.get("/health")
async def health_check():
    """
    Health check endpoint.
    """
    return {"status": "healthy"}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
