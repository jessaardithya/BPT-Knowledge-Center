import fitz
import logging
from typing import List, Dict, Any
from sentence_transformers import SentenceTransformer
from models.schemas import ContentItem
from core.config import settings

logger = logging.getLogger(__name__)


class DocumentProcessor:
    def __init__(self):
        logger.info(f"Loading Embedding Model: {settings.EMBEDDING_MODEL}")
        self.model = SentenceTransformer(settings.EMBEDDING_MODEL)
        logger.info("Model loaded")

    def embed_text(self, text: str) -> List[float]:
        # Simple wrapper around model.encode
        return self.model.encode(text).tolist()

    def process_pdf(self, file_path: str, filename: str) -> List[ContentItem]:
        logger.info(f"Processing file: {filename}")
        doc = fitz.open(file_path)
        extracted_data = []

        try:
            for page_num, page in enumerate(doc):
                text = page.get_text()
                if not text.strip():
                    continue

                blocks = text.split('\n\n')
                for i, block in enumerate(blocks):
                    clean_text = block.strip()
                    if len(clean_text) < 20:
                        continue

                    # Embedding generation is blocking, so this function is blocking
                    embedding = self.embed_text(clean_text)

                    extracted_data.append(ContentItem(
                        element_id=f"p{page_num}_b{i}",
                        text=clean_text,
                        type="text",
                        metadata={"page": page_num + 1, "source": filename},
                        vector=embedding
                    ))
        finally:
            doc.close()

        logger.info(f"Extracted {len(extracted_data)} items from {filename}")
        return extracted_data


# Singleton instance
document_processor = DocumentProcessor()
