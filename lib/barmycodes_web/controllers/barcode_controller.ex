defmodule BarmycodesWeb.BarcodeController do
  use BarmycodesWeb, :controller

  alias Barmycodes.Barcodes

  def index(conn, %{ "type" => barcodes_type, "b" => barcode_values }) do
    barcodes =
      Enum.uniq(barcode_values)
      |> Enum.reject(& &1 == nil || String.trim(&1) == "")
      |> Enum.map(& Barcodes.generate!(barcodes_type, &1))

    render(conn, "index.html", conn: conn, barcodes: barcodes)
  end
end
