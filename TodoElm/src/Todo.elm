module Todo exposing (..)

import Html exposing (..)
import Browser exposing (sandbox)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput)


-- MAIN


main : Program () Model Msg
main =
    Browser.sandbox { init = initialModel, update = update, view = view }
    
-- MODEL


type alias Todo =
    { id : Int
    , title : String
    , status : Bool
    }


type alias Model =
    { todos : List Todo
    , id : Int
    , title : String
    }


initialModel : Model
initialModel =
    { todos = []
    , id = 0
    , title = ""
    }


newTodoItem : String -> Int -> Todo
newTodoItem title id =
    { title = title
    , id = id
    , status = False
    }


-- UPDATE

type Msg
    = AddTodoItem
    | UpdateTitle String


update : Msg -> Model -> Model
update msg model =
    case msg of
        AddTodoItem ->
            { model
                | id = model.id + 1
                , title = ""
                , todos =
                    if String.isEmpty model.title then
                        model.todos

                    else
                        model.todos ++ [ newTodoItem model.title model.id ]
            }

        UpdateTitle value ->
            { model | title = value }




-- VIEW


view : Model -> Html Msg
view model =
    div [ class "container" ]
        [ viewHero
        , viewAddTodoInput model.title
        , viewTodoList model.todos
        ]


viewHero : Html Msg
viewHero =
    div [ class "hero" ] [ text "A Todo List Built In Elm" ]


viewAddTodoInput : String -> Html Msg
viewAddTodoInput title =
    div [ class "item-group" ]
        [ div [ class "item-title" ] [ text "Add a todo item" ]
        , input
            [ value title
            , onInput UpdateTitle
            , placeholder "Type Here"
            , class "item-input"
            ]
            []
        , button [ onClick AddTodoItem, class "item-button" ] [ text "add item" ]
        ]



-- VIEW TODOS


viewTodoList : List Todo -> Html Msg
viewTodoList todos =
    let
        todoList =
            List.map viewTodoItem todos
    in
    div [ class "item-group" ]
        [ div [ class "item-title" ] [ text "Things to do" ]
        , ul [ class "todo-item-group" ] todoList
        ]


viewTodoItem : Todo -> Html Msg
viewTodoItem todo =
    if todo.status == False then
        li [ class "todo-item" ]
            [ span [ class "todo-item-title" ]
                [ text todo.title ]
  
            ]
    else
        text ""




