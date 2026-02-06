from typing import List, Dict, Any
from pydantic import BaseModel


class ContentItem(BaseModel):
    element_id: str
    text: str
    type: str
    metadata: Dict[str, Any]
    vector: List[float]


class ParseResponse(BaseModel):
    filename: str
    content_type: str
    element_count: int
    data: List[ContentItem]


class EmbedRequest(BaseModel):
    text: str


class EmbedResponse(BaseModel):
    text: str
    vector: List[float]
