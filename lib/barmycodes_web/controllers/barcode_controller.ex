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

  def png(conn, %{ "type" => barcode_type, "b" => barcode_value }) do
    barcode_value = List.first(barcode_value)
    barcode = Barcodes.generate!(barcode_type, barcode_value)

    options = [
      filename: "barmycodes_#{barcode_value}.png",
    ]
    send_download(conn, {:binary, barcode.image}, options)
  end
end
