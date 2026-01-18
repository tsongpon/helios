package co.abctech.helios.controller;

import co.abctech.helios.model.Statement;
import co.abctech.helios.service.StatementService;
import java.io.IOException;
import org.apache.pdfbox.Loader;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.text.PDFTextStripper;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/v1/statements")
public class StatementController {

    private final StatementService statementService;

    public StatementController(StatementService statementService) {
        this.statementService = statementService;
    }

    @PostMapping(
        path = "/extract",
        consumes = MediaType.MULTIPART_FORM_DATA_VALUE,
        produces = MediaType.APPLICATION_JSON_VALUE
    )
    public ResponseEntity<?> extractStatement(
        @RequestParam("file") MultipartFile file,
        @RequestParam(value = "password", required = false) String password
    ) {
        if (file.isEmpty()) {
            return ResponseEntity.badRequest().body(
                new ErrorResponse("No file uploaded")
            );
        }

        if (!isPdfFile(file)) {
            return ResponseEntity.badRequest().body(
                new ErrorResponse("Uploaded file is not a PDF")
            );
        }

        try {
            // Step 1: Extract text from PDF
            String extractedText = extractPdfText(file, password);

            // Step 2: Parse statement using LLM service
            Statement statement = statementService.parseStatement(
                extractedText
            );

            // Step 3: Return parsed statement model
            return ResponseEntity.ok(statement);
        } catch (IOException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(
                new ErrorResponse("Error processing PDF: " + e.getMessage())
            );
        } catch (IllegalArgumentException e) {
            return ResponseEntity.badRequest().body(
                new ErrorResponse(e.getMessage())
            );
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(
                new ErrorResponse("Unexpected error: " + e.getMessage())
            );
        }
    }

    @PostMapping(
        path = "/extract-text",
        consumes = MediaType.MULTIPART_FORM_DATA_VALUE,
        produces = MediaType.TEXT_PLAIN_VALUE
    )
    public ResponseEntity<String> extractTextFromPdf(
        @RequestParam("file") MultipartFile file,
        @RequestParam(value = "password", required = false) String password
    ) {
        if (file.isEmpty()) {
            return ResponseEntity.badRequest().body("No file uploaded");
        }

        if (!isPdfFile(file)) {
            return ResponseEntity.badRequest().body(
                "Uploaded file is not a PDF"
            );
        }

        try {
            String extractedText = extractPdfText(file, password);
            return ResponseEntity.ok(extractedText);
        } catch (IOException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(
                "Error processing PDF: " + e.getMessage()
            );
        } catch (Exception e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(
                "Unexpected error: " + e.getMessage()
            );
        }
    }

    private boolean isPdfFile(MultipartFile file) {
        String contentType = file.getContentType();
        String filename = file.getOriginalFilename();

        return (
            (contentType != null && contentType.equals("application/pdf")) ||
            (filename != null && filename.toLowerCase().endsWith(".pdf"))
        );
    }

    private String extractPdfText(MultipartFile file, String password)
        throws IOException {
        try (
            PDDocument document = Loader.loadPDF(
                file.getBytes(),
                password != null ? password : ""
            )
        ) {
            if (document.isEncrypted() && password == null) {
                throw new IOException(
                    "PDF is password protected but no password was provided"
                );
            }

            PDFTextStripper stripper = new PDFTextStripper();
            return stripper.getText(document);
        }
    }

    private record ErrorResponse(String message) {}
}
