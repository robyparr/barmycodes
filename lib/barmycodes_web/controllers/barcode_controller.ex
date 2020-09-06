defmodule BarmycodesWeb.BarcodeController do
  use BarmycodesWeb, :controller

  alias Barmycodes.Barcodes

  def index(conn, %{ "b" => barcode_values }) do
    barcodes =
      Enum.uniq(barcode_values)
      |> Enum.map(&(Barcodes.generate!("code128", &1)))

    render(conn, "index.html", conn: conn, barcodes: barcodes)
  end
end
