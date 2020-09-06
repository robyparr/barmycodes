defmodule BarmycodesWeb.PageController do
  use BarmycodesWeb, :controller

  def index(conn, _params) do
    render(conn, "index.html")
  end
end
