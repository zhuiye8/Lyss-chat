"""
File management and document processing.

This module defines models for handling user file uploads,
document processing status, and vector storage integration.
"""

import uuid
from typing import Optional

from sqlalchemy import BigInteger, ForeignKey, Integer, String, Text
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.db.base import Base


class UploadedFile(Base):
    """
    User uploaded files for document processing and Q&A.
    
    Tracks file uploads, processing status, and vector storage
    for document-based AI interactions.
    """
    
    __tablename__ = "uploaded_file"
    
    # User relationship
    user_id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("user.id", ondelete="CASCADE"),
        nullable=False,
        index=True
    )
    
    # File basic information
    original_filename: Mapped[str] = mapped_column(String(255), nullable=False)
    file_path: Mapped[str] = mapped_column(String(500), nullable=False, unique=True)
    file_size: Mapped[int] = mapped_column(BigInteger, nullable=False)
    mime_type: Mapped[str] = mapped_column(String(100), nullable=False)
    
    # Processing status
    status: Mapped[str] = mapped_column(
        String(50),
        default="pending",
        nullable=False,
        index=True,
        comment="pending, processing, completed, failed"
    )
    processing_error: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    
    # Vector storage information
    total_chunks: Mapped[int] = mapped_column(Integer, default=0, nullable=False)
    qdrant_collection: Mapped[Optional[str]] = mapped_column(
        String(100),
        nullable=True,
        comment="Qdrant collection name for vector storage"
    )
    
    # Processing metadata
    extracted_text_length: Mapped[Optional[int]] = mapped_column(
        Integer,
        nullable=True,
        comment="Length of extracted text content"
    )
    processing_duration_ms: Mapped[Optional[int]] = mapped_column(
        Integer,
        nullable=True,
        comment="Processing duration in milliseconds"
    )
    
    # Relationships
    user = relationship("User", back_populates="uploaded_files")
    
    @property
    def is_pending(self) -> bool:
        """Check if file is pending processing."""
        return self.status == "pending"
    
    @property
    def is_processing(self) -> bool:
        """Check if file is currently being processed."""
        return self.status == "processing"
    
    @property
    def is_completed(self) -> bool:
        """Check if file processing is completed."""
        return self.status == "completed"
    
    @property
    def is_failed(self) -> bool:
        """Check if file processing failed."""
        return self.status == "failed"
    
    @property
    def file_size_mb(self) -> float:
        """Get file size in megabytes."""
        return self.file_size / (1024 * 1024)
    
    @property
    def file_extension(self) -> str:
        """Get file extension from filename."""
        return self.original_filename.split(".")[-1].lower() if "." in self.original_filename else ""
    
    @property
    def is_vectorized(self) -> bool:
        """Check if file has been vectorized and stored in Qdrant."""
        return self.qdrant_collection is not None and self.total_chunks > 0
    
    def mark_processing(self) -> None:
        """Mark file as currently processing."""
        self.status = "processing"
        self.processing_error = None
    
    def mark_completed(self, total_chunks: int, qdrant_collection: str) -> None:
        """
        Mark file processing as completed.
        
        Args:
            total_chunks: Number of text chunks created
            qdrant_collection: Qdrant collection name where vectors are stored
        """
        self.status = "completed"
        self.total_chunks = total_chunks
        self.qdrant_collection = qdrant_collection
        self.processing_error = None
    
    def mark_failed(self, error_message: str) -> None:
        """
        Mark file processing as failed.
        
        Args:
            error_message: Error description
        """
        self.status = "failed"
        self.processing_error = error_message
    
    def __repr__(self) -> str:
        return (
            f"<UploadedFile(id={self.id}, filename={self.original_filename}, "
            f"status={self.status}, user_id={self.user_id})>"
        )