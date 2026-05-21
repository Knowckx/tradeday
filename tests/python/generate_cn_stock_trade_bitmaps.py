import json
from datetime import date
from pathlib import Path
import sys


ROOT = Path(__file__).resolve().parent.parent.parent
TRUTH_TABLE_PATH = ROOT / "tests" / "testdata" / "cn_stock_truth_table.json"
OUTPUT_PATH = ROOT / "core" / "data" / "cn_days_bitmaps.go"


def load_truth_table(path: Path) -> tuple[date, date, dict[str, bool]]:
    payload = json.loads(path.read_text(encoding="utf-8"))

    if payload.get("calendar_id") != "cn_stock":
        raise ValueError(f"unexpected calendar_id: {payload.get('calendar_id')}")

    start = date.fromisoformat(payload["start"])
    end = date.fromisoformat(payload["end"])
    days = payload["days"]
    if not isinstance(days, dict):
        raise ValueError("truth table days must be a JSON object")

    return start, end, days


def build_year_bitmap(year: int, days: dict[str, bool]) -> list[int]:
    bits = [0] * 6
    for day, is_trade_day in days.items():
        if not is_trade_day:
            continue

        current_day = date.fromisoformat(day)
        if current_day.year != year:
            continue

        bit_index = current_day.timetuple().tm_yday - 1
        word_index = bit_index // 64
        bit_offset = bit_index % 64
        if word_index >= len(bits):
            raise ValueError(f"day out of bitmap range: {day}")

        bits[word_index] |= 1 << bit_offset

    return bits


def render_go_source(start: date, end: date, days: dict[str, bool]) -> str:
    years = list(range(start.year, end.year + 1))
    lines: list[str] = [
        "package data",
        "",
        "// CNStockTradeBitmaps 保存中国 A 股各年的交易日真值位图。",
        "var CNStockTradeBitmaps = YearTradeBitmaps{",
    ]

    for year in years:
        bits = build_year_bitmap(year, days)
        lines.append(f"\t{year}: {{")
        lines.append(f"\t\tYear: {year},")
        lines.append("\t\tBits: [6]uint64{")
        for value in bits:
            lines.append(f"\t\t\t{value},")
        lines.append("\t\t},")
        lines.append("\t},")

    lines.extend(
        [
            "}",
            "",
            "func init() {",
            "\tCNStockTradeBitmaps.mustAlignYearKeys()",
            "}",
            "",
        ]
    )

    return "\n".join(lines)


def main() -> int:
    if len(sys.argv) != 1:
        raise SystemExit("usage: python generate_cn_stock_trade_bitmaps.py")

    start, end, days = load_truth_table(TRUTH_TABLE_PATH)
    source = render_go_source(start, end, days)
    OUTPUT_PATH.write_text(source, encoding="utf-8")
    print(OUTPUT_PATH)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
