defmodule ProjectOffer.AccountsFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `ProjectOffer.Accounts` context.
  """
alias ProjectOffer.Accounts.User

  def unique_user_email, do: "user#{System.unique_integer()}@example.com"
  def valid_user_password, do: "hello world!"

  def valid_user_attributes(attrs \\ %{}) do
    Enum.into(attrs, %{
      email: unique_user_email(),
      password: valid_user_password()
    })
  end

  def user_fixture(attrs \\ %{}) do
    {:ok, user} =
      attrs
      |> valid_user_attributes()
      |> ProjectOffer.Accounts.register_user()

    {:ok, user} = confirm_user(user)

    user
  end

  def user_fixture_unconfirmed(attrs \\ %{}) do
    {:ok, user} =
      attrs
      |> valid_user_attributes()
      |> ProjectOffer.Accounts.register_user()

    user
  end

  def extract_user_token(fun) do
    {:ok, captured_email} = fun.(&"[TOKEN]#{&1}[TOKEN]")
    [_, token | _] = String.split(captured_email.text_body, "[TOKEN]")
    token
  end

  defp confirm_user(user) do
    alias ProjectOffer.Repo

    Repo.update(User.confirm_changeset(user))
  end
end
