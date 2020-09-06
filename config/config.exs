# This file is responsible for configuring your application
# and its dependencies with the aid of the Mix.Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.

# General application configuration
use Mix.Config

config :barmycodes,
  ecto_repos: [Barmycodes.Repo]

# Configures the endpoint
config :barmycodes, BarmycodesWeb.Endpoint,
  url: [host: "localhost"],
  secret_key_base: "L7m6RSvmvlrCrDfTKYaj5exwwHV3+cwqJW0rc+yign+65PgC7df3ISWh4ZmhvYXL",
  render_errors: [view: BarmycodesWeb.ErrorView, accepts: ~w(html json), layout: false],
  pubsub_server: Barmycodes.PubSub,
  live_view: [signing_salt: "Cia5KISD"]

# Configures Elixir's Logger
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:request_id]

# Use Jason for JSON parsing in Phoenix
config :phoenix, :json_library, Jason

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{Mix.env()}.exs"
