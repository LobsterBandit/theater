import { useEffect, useReducer, useState } from "react";

const initialState = {
  loading: true,
  plexWebhooks: [],
  error: null,
};

function reducer(state, action) {
  switch (action.type) {
    case "FETCH_ERROR":
      return {
        ...state,
        loading: false,
        error: action.payload,
      };
    case "FETCH_START":
      return {
        ...state,
        loading: true,
        error: null,
      };
    case "FETCH_SUCCESS":
      return {
        loading: false,
        plexWebhooks: action.payload,
        error: null,
      };
    default:
      return state;
  }
}

export function usePlexWebhooks() {
  const [shouldFetch, setShouldFetch] = useState(true);
  const [state, dispatch] = useReducer(reducer, initialState);

  const refetch = () => setShouldFetch(true);

  useEffect(() => {
    async function fetchPlexWebhooks() {
      dispatch({ type: "FETCH_START" });
      try {
        const resp = await fetch("/plex");
        const respData = await resp.json();
        if (respData) {
          dispatch({ type: "FETCH_SUCCESS", payload: respData });
        }
      } catch (error) {
        console.error(error);
        dispatch({ type: "FETCH_ERROR", payload: error });
      } finally {
        setShouldFetch(false);
      }
    }

    shouldFetch && fetchPlexWebhooks();
  }, [shouldFetch]);

  return [state, refetch];
}
