defmodule Barmycodes.Barcodes do
  alias Barmycodes.Barcodes.Barcode

  def generate!(type, value) do
    {:ok, barcode} = generate(type, value)
    barcode
  end

  def generate("code128", value) do
    file_path = temp_generation_path()
    Barlix.Code128.encode!(value)
    |> Barlix.PNG.print(file: file_path, xdim: 3, margin: 2)
    add_value_label_to_barcode_image(file_path, value)

    encoded_image = :base64.encode(File.read!(file_path))
    File.rm! file_path

    barcode = %Barcode{type: "code128", value: value, encoded_image: encoded_image}
    {:ok, barcode}
  end

  def generate("qr_code", value) do
    encoded_barcode =
      EQRCode.encode(value)
      |> EQRCode.png()
      |> :base64.encode()

    barcode = %Barcode{type: "qr", value: value, encoded_image: encoded_barcode}
    {:ok, barcode}
  end

  defp temp_generation_path do
    "#{System.tmp_dir()}/#{Ecto.UUID.generate}"
  end

  defp add_value_label_to_barcode_image(file_path, value_label) do
    image =
      Mogrify.open(file_path)
      |> Mogrify.verbose

    image
    |> Mogrify.custom("gravity", "North")
    |> Mogrify.custom("extent", "#{image.width}x#{image.height + 40}")
    |> Mogrify.custom("pointsize", 40)
    |> Mogrify.custom("gravity", "South")
    |> Mogrify.custom("annotate", "+0+0 #{value_label}")
    |> Mogrify.save(path: file_path)
  end
end
