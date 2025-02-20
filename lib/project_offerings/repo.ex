defmodule ProjectOfferings.Repo do
  use Ecto.Repo,
    otp_app: :project_offerings,
    adapter: Ecto.Adapters.Postgres
end
