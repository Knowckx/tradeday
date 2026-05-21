import json
from pathlib import Path
import sys
from datetime import date, timedelta

import exchange_calendars as xc


TRUTH_TABLE_START = date(2015, 1, 1)
TRUTH_TABLE_END = date(2026, 12, 31)
OUTPUT_PATH = Path(__file__).resolve().parent.parent / "testdata" / "us_stock_truth_table.json"
US_STOCK_CALENDAR_NAME = "XNYS"


def build_truth_table(start: date, end: date) -> dict[str, bool]:
    calendar = xc.get_calendar(US_STOCK_CALENDAR_NAME)
    truth_table: dict[str, bool] = {}

    current = start
    while current <= end:
        day = current.isoformat()
        truth_table[day] = is_trade_day(calendar, current)
        current += timedelta(days=1)

    return truth_table


def is_trade_day(calendar, day: date) -> bool:
    return calendar.is_session(day.isoformat())


def build_payload() -> dict:
    return {
        "calendar_id": "us_stock",
        "start": TRUTH_TABLE_START.isoformat(),
        "end": TRUTH_TABLE_END.isoformat(),
        "generated_by": "tests/python/generate_us_stock_truth_table.py",
        "sources": {
            "2015-2026": "exchange_calendars XNYS 4.11.3",
        },
        "days": build_truth_table(TRUTH_TABLE_START, TRUTH_TABLE_END),
    }


def main() -> int:
    if len(sys.argv) != 1:
        raise SystemExit("usage: python generate_us_stock_truth_table.py")

    payload = build_payload()
    OUTPUT_PATH.parent.mkdir(parents=True, exist_ok=True)
    OUTPUT_PATH.write_text(
        json.dumps(payload, ensure_ascii=True, sort_keys=True, indent=2) + "\n",
        encoding="utf-8",
    )
    print(OUTPUT_PATH)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
