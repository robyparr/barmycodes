defmodule Barmycodes.Barcodes.Barcode do
  defstruct [
    :type,
    :width,
    :height,
    :value,
    :image,
    :encoded_image,
    :image_path,
  ]
end
