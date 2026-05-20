import json
from pathlib import Path
import sys
from datetime import date, timedelta

import exchange_calendars as xc


TRUTH_TABLE_START = date(2015, 1, 1)
TRUTH_TABLE_END = date(2026, 12, 31)
EXCHANGE_CALENDARS_END = date(2025, 12, 31)
OUTPUT_PATH = Path(__file__).resolve().parent.parent / "testdata" / "cn_stock_truth_table.json"
CN_STOCK_CLOSED_DAYS_2026 = {
    "2026-01-01",
    "2026-01-02",
    "2026-02-16",
    "2026-02-17",
    "2026-02-18",
    "2026-02-19",
    "2026-02-20",
    "2026-02-23",
    "2026-04-06",
    "2026-05-01",
    "2026-05-04",
    "2026-05-05",
    "2026-06-19",
    "2026-09-25",
    "2026-10-01",
    "2026-10-02",
    "2026-10-05",
    "2026-10-06",
    "2026-10-07",
}


def build_truth_table(start: date, end: date) -> dict[str, bool]:
    calendar = xc.get_calendar("XSHG")
    truth_table: dict[str, bool] = {}

    current = start
    while current <= end:
        day = current.isoformat()
        truth_table[day] = is_trade_day(calendar, current)
        current += timedelta(days=1)

    return truth_table


def is_trade_day(calendar, day: date) -> bool:
    if day <= EXCHANGE_CALENDARS_END:
        return calendar.is_session(day.isoformat())

    if day.year != 2026:
        raise ValueError(f"unsupported day: {day.isoformat()}")

    if day.weekday() >= 5:
        return False

    return day.isoformat() not in CN_STOCK_CLOSED_DAYS_2026


def build_payload() -> dict:
    return {
        "calendar_id": "cn_stock",
        "start": TRUTH_TABLE_START.isoformat(),
        "end": TRUTH_TABLE_END.isoformat(),
        "generated_by": "tests/python/generate_cn_stock_truth_table.py",
        "sources": {
            "2015-2025": "exchange_calendars XSHG 4.11.3",
            "2026": "https://www.sse.com.cn/disclosure/announcement/general/c/c_20251222_10802507.shtml",
        },
        "days": build_truth_table(TRUTH_TABLE_START, TRUTH_TABLE_END),
    }


def main() -> int:
    if len(sys.argv) != 1:
        raise SystemExit("usage: python generate_cn_stock_truth_table.py")

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
