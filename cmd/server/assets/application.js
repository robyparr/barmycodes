const barcodeInput = document.querySelector("#barcode-values");

window.addEventListener("load", function () {
  const urlSearchParams = new URLSearchParams(window.location.search);
  barcodeInput.value = urlSearchParams.getAll("b[]").join("\n");
});

function generateBarcodes() {
  const barcodeInputValue = barcodeInput.value.trim();
  if (barcodeInputValue.length == 0) return;

  const barcodeType = document.querySelector('input[name="type"]:checked');
  const baseURL = window.location.protocol + "//" + window.location.host + "/";

  const barcodeValues = barcodeInputValue.split("\n");
  const params = ["type=" + barcodeType.value];

  for (var i = 0; i < barcodeValues.length; i++) {
    params.push("b[]=" + encodeURIComponent(barcodeValues[i]));
  }

  window.location.href = baseURL + "?" + params.join("&");
}

document
  .querySelector("#generate-barcodes")
  .addEventListener("click", generateBarcodes);

barcodeInput.addEventListener("keydown", function (e) {
  if ((e.ctrlKey || e.metaKey) && (e.keyCode === 13 || e.keyCode === 10)) {
    generateBarcodes();
  }
});
