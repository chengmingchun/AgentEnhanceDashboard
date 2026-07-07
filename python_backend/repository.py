from __future__ import annotations

from typing import Protocol

from .database import SQLiteConnectionFactory
from .models import Record


class RecordRepository(Protocol):
    def list_records(self) -> list[Record]:
        """Return all dashboard records in display order."""


class SQLiteRecordRepository:
    def __init__(self, connection_factory: SQLiteConnectionFactory) -> None:
        self._connection_factory = connection_factory

    def list_records(self) -> list[Record]:
        with self._connection_factory.connect() as connection:
            rows = connection.execute(
                """
                SELECT
                    id,
                    requirement,
                    mr,
                    mr_url,
                    owner,
                    group_name,
                    lines,
                    efficiency,
                    score,
                    problem,
                    record_date,
                    status
                FROM records
                ORDER BY record_date DESC, id DESC
                """
            ).fetchall()

        return [
            Record(
                id=row["id"],
                requirement=row["requirement"],
                mr=row["mr"],
                mr_url=row["mr_url"],
                owner=row["owner"],
                group=row["group_name"],
                lines=row["lines"],
                efficiency=row["efficiency"],
                score=row["score"],
                problem=row["problem"],
                date=row["record_date"],
                status=row["status"],
            )
            for row in rows
        ]
