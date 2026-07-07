from __future__ import annotations

import sqlite3
from contextlib import contextmanager
from pathlib import Path
from typing import Iterator

from .seed import SEED_RECORDS


class SQLiteConnectionFactory:
    def __init__(self, database_path: Path) -> None:
        self._database_path = database_path

    @contextmanager
    def connect(self) -> Iterator[sqlite3.Connection]:
        self._database_path.parent.mkdir(parents=True, exist_ok=True)
        connection = sqlite3.connect(self._database_path)
        connection.row_factory = sqlite3.Row
        try:
            yield connection
            connection.commit()
        except Exception:
            connection.rollback()
            raise
        finally:
            connection.close()


class DatabaseInitializer:
    def __init__(self, connection_factory: SQLiteConnectionFactory) -> None:
        self._connection_factory = connection_factory

    def initialize(self) -> None:
        with self._connection_factory.connect() as connection:
            self._migrate(connection)
            self._seed_if_empty(connection)

    def _migrate(self, connection: sqlite3.Connection) -> None:
        connection.executescript(
            """
            CREATE TABLE IF NOT EXISTS records (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                requirement TEXT NOT NULL,
                mr TEXT NOT NULL,
                mr_url TEXT NOT NULL,
                owner TEXT NOT NULL,
                group_name TEXT NOT NULL,
                lines INTEGER NOT NULL,
                efficiency INTEGER NOT NULL,
                score REAL NOT NULL,
                problem TEXT NOT NULL,
                record_date TEXT NOT NULL,
                status TEXT NOT NULL,
                created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
            );

            CREATE INDEX IF NOT EXISTS idx_records_date ON records(record_date);
            CREATE INDEX IF NOT EXISTS idx_records_owner ON records(owner);
            CREATE INDEX IF NOT EXISTS idx_records_group ON records(group_name);
            """
        )

    def _seed_if_empty(self, connection: sqlite3.Connection) -> None:
        count = connection.execute("SELECT COUNT(*) FROM records").fetchone()[0]
        if count:
            return

        connection.executemany(
            """
            INSERT INTO records (
                requirement, mr, mr_url, owner, group_name, lines,
                efficiency, score, problem, record_date, status
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            [
                (
                    item.requirement,
                    item.mr,
                    item.mr_url,
                    item.owner,
                    item.group,
                    item.lines,
                    item.efficiency,
                    item.score,
                    item.problem,
                    item.date,
                    item.status,
                )
                for item in SEED_RECORDS
            ],
        )
