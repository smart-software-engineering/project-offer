defmodule ProjectOffer.Repo do
  use Ecto.Repo,
    otp_app: :project_offer,
    adapter: Ecto.Adapters.Postgres
end
