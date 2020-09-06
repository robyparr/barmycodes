defmodule Barmycodes.Barcodes do
  alias Barmycodes.Barcodes.Barcode

  def generate!(type, value) do
    {:ok, barcode} = generate(type, value)
    barcode
  end

  def generate("code128", value) do
    with_tempfile(fn file_path ->
      Barlix.Code128.encode!(value)
      |> Barlix.PNG.print(file: file_path, xdim: 3, margin: 2)
      {width, height} = add_value_label_to_barcode_image(file_path, value)

      image = File.read!(file_path)
      barcode =
        %Barcode{
          type: "code128",
          width: width,
          height: height,
          value: value,
          image: image,
          encoded_image: :base64.encode(image),
          image_path: file_path,
        }
      {:ok, barcode}
    end)
  end

  def generate("qr_code", value) do
    with_tempfile(fn file_path ->
      barcode =
        EQRCode.encode(value)
        |> EQRCode.png()

      File.write!(file_path, barcode, [:binary])

      barcode =
        %Barcode{
          type: "qr",
          width: 275,
          height: 275,
          value: value,
          image: barcode,
          encoded_image: :base64.encode(barcode),
          image_path: file_path,
        }
      {:ok, barcode}
    end)
  end

  defp add_value_label_to_barcode_image(file_path, value_label) do
    image =
      Mogrify.open(file_path)
      |> Mogrify.verbose

    new_height = image.height + 40

    image
    |> Mogrify.custom("gravity", "North")
    |> Mogrify.custom("extent", "#{image.width}x#{new_height}")
    |> Mogrify.custom("pointsize", 40)
    |> Mogrify.custom("gravity", "South")
    |> Mogrify.custom("annotate", "+0+0 #{value_label}")
    |> Mogrify.save(path: file_path)

    {image.width, new_height}
  end

  defp with_tempfile(on_create) do
    {:ok, path} = Briefly.create
    on_create.(path)
  end
end
