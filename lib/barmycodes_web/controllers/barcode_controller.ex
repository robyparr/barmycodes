defmodule BarmycodesWeb.BarcodeController do
  use BarmycodesWeb, :controller

  alias Barmycodes.Barcodes

  def index(conn, %{ "type" => barcodes_type, "b" => barcode_values }) do
    barcodes =
      scrub_barcode_values(barcode_values)
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

  def pdf(conn, %{ "type" => barcodes_type, "b" => barcode_values }) do
    barcodes =
      scrub_barcode_values(barcode_values)
      |> Enum.map(& Barcodes.generate!(barcodes_type, &1))

    first_barcode = List.first(barcodes)
    {:ok, pdf} = Pdf.new([size: [first_barcode.width, first_barcode.height], compress: true])

    barcodes
    |> Enum.with_index
    |> Enum.each(fn {barcode, index} ->
      if index > 0, do: Pdf.add_page(pdf, [barcode.width, barcode.height])
      Pdf.add_image(pdf, {0, 0}, barcode.image_path)
    end)

    options =[
      filename: "barmycodes.pdf",
      disposition: :inline,
    ]
    send_download(conn, {:binary, Pdf.export(pdf)}, options)
  end

  defp scrub_barcode_values(barcode_values) do
    Enum.uniq(barcode_values)
    |> Enum.reject(& &1 == nil || String.trim(&1) == "")
  end
end
