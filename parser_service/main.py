import os
import shutil
import tempfile
import asyncio
from fastapi import FastAPI, UploadFile, File, HTTPException
from concurrent.futures import ThreadPoolExecutor

from core.config import settings
from core.logging_config import setup_logging
from models.schemas import ParseResponse, EmbedRequest, EmbedResponse
from services.document_processor import document_processor

# Setup logging
logger = setup_logging()

app = FastAPI(title=settings.PROJECT_NAME)

# Create a thread pool for CPU-bound tasks
thread_pool = ThreadPoolExecutor(max_workers=os.cpu_count() or 1)


@app.post(f"{settings.API_V1_STR}/parse", response_model=ParseResponse)
async def parse_document(file: UploadFile = File(...)):
    filename = file.filename
    logger.info(f"Received file upload: {filename}")

    # Create a temporary file to save the uploaded content
    try:
        with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf") as tmp_file:
            shutil.copyfileobj(file.file, tmp_file)
            temp_path = tmp_file.name
    except Exception as e:
        logger.error(f"Failed to save temporary file: {e}")
        raise HTTPException(status_code=500, detail="Failed to upload file")

    try:
        # Run processing in a separate thread to avoid blocking the event loop
        loop = asyncio.get_event_loop()
        extracted_data = await loop.run_in_executor(
            thread_pool,
            document_processor.process_pdf,
            temp_path,
            filename
        )

        return ParseResponse(
            filename=filename,
            content_type="application/pdf",
            element_count=len(extracted_data),
            data=extracted_data
        )

    except Exception as e:
        logger.error(
            f"Error processing document {filename}: {e}", exc_info=True)
        raise HTTPException(
            status_code=500, detail=f"Processing failed: {str(e)}")

    finally:
        # Clean up the temporary file
        if os.path.exists(temp_path):
            try:
                os.remove(temp_path)
            except OSError as e:
                logger.warning(f"Failed to remove temp file {temp_path}: {e}")

@app.post(f"{settings.API_V1_STR}/embed", response_model=EmbedResponse)
async def embed_text(req: EmbedRequest):
    try:
        # Run processing in a separate thread to avoid blocking the event loop
        loop = asyncio.get_event_loop()
        embedding = await loop.run_in_executor(
            thread_pool,
            document_processor.embed_text,
            req.text
        )

        return EmbedResponse(
            text=req.text,
            vector=embedding
        )
    except Exception as e:
        logger.error(f"Error embedding text: {e}", exc_info=True)
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host=settings.HOST,
                port=settings.PORT, reload=True)
