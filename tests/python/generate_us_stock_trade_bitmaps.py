import sys

from generate_trade_bitmaps import BITMAP_SPECS, write_bitmap_file


def main() -> int:
    if len(sys.argv) != 1:
        raise SystemExit("usage: python generate_us_stock_trade_bitmaps.py")

    output_path = write_bitmap_file(BITMAP_SPECS["us_stock"])
    print(output_path)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
