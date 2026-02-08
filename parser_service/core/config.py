import os
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    API_V1_STR: str = "/api/v1"
    PROJECT_NAME: str = "Parser Service"
    EMBEDDING_MODEL: str = "all-MiniLM-L6-v2"
    HOST: str = "0.0.0.0"
    PORT: int = 8000
    HF_TOKEN: str = ""

    class Config:
        env_file = ".env"


settings = Settings()

