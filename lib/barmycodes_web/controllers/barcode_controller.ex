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
      |> List.first()

    pdf =
      Pdf.build([size: [barcodes.width, barcodes.height], compress: true], fn pdf ->
        pdf
        |> Pdf.add_image({0, 0}, barcodes.image_path)
        |> Pdf.export()
      end)

    options =[
      filename: "barmycodes.pdf",
      disposition: :inline,
    ]
    send_download(conn, {:binary, pdf}, options)
  end

  defp scrub_barcode_values(barcode_values) do
    Enum.uniq(barcode_values)
    |> Enum.reject(& &1 == nil || String.trim(&1) == "")
  end
end
