defmodule Barmycodes.Repo do
  use Ecto.Repo,
    otp_app: :barmycodes,
    adapter: Ecto.Adapters.Postgres
end
