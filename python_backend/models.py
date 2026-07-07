from __future__ import annotations

from dataclasses import dataclass
from typing import Any


@dataclass(frozen=True)
class Record:
    id: int | None
    requirement: str
    mr: str
    mr_url: str
    owner: str
    group: str
    lines: int
    efficiency: int
    score: float
    problem: str
    date: str
    status: str

    def to_api(self) -> dict[str, Any]:
        return {
            "id": self.id,
            "requirement": self.requirement,
            "mr": self.mr,
            "mrUrl": self.mr_url,
            "owner": self.owner,
            "group": self.group,
            "lines": self.lines,
            "efficiency": self.efficiency,
            "score": self.score,
            "problem": self.problem,
            "date": self.date,
            "status": self.status,
        }
