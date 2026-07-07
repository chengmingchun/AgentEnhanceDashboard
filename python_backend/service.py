from __future__ import annotations

from .models import Record
from .repository import RecordRepository


class RecordService:
    def __init__(self, repository: RecordRepository) -> None:
        self._repository = repository

    def list_records(self) -> list[Record]:
        return self._repository.list_records()

    def list_records_for_api(self) -> list[dict]:
        return [record.to_api() for record in self.list_records()]
