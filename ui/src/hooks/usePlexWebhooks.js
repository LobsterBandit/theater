import { useCallback, useEffect, useReducer } from "react";

const initialState = {
  loading: true,
  plexWebhooks: [],
  total: 0,
  error: null,
  pagination: {
    pageIndex: 0,
    pageSize: 10,
    orderBy: "id",
    sortBy: "asc",
  },
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
        total: -1,
        error: null,
      };
    case "UPDATE_PAGINATION":
      return {
        ...state,
        pagination: {
          ...state.pagination,
          ...action.payload,
        },
      };
    default:
      return state;
  }
}

function buildURL({ pageIndex, pageSize, orderBy, sortBy }) {
  const baseURL = "/plex";
  return pageIndex >= 0 && pageSize && orderBy && sortBy
    ? `${baseURL}?limit=${pageSize}&offset=${
        pageIndex * pageSize
      }&orderBy=${orderBy}&sortBy=${sortBy}`
    : baseURL;
}

export function usePlexWebhooks({
  fetchOnMount = false,
  pageIndex,
  pageSize,
  orderBy,
  sortBy,
} = {}) {
  const [state, dispatch] = useReducer(
    reducer,
    {
      ...initialState,
      pagination: {
        ...initialState.pagination,
        pageIndex,
        pageSize,
        orderBy,
        sortBy,
      },
      fetchOnMount,
    },
    (args) => ({
      ...args,
      loading: args.fetchOnMount,
      pagination: {
        ...args.pagination,
      },
    })
  );

  const fetchPlexWebhooks = useCallback(
    async ({ pageIndex, pageSize, orderBy, sortBy }) => {
      dispatch({ type: "FETCH_START" });
      try {
        const url = buildURL({ pageIndex, pageSize, orderBy, sortBy });
        const resp = await fetch(url);
        const respData = await resp.json();
        if (respData) {
          dispatch({ type: "FETCH_SUCCESS", payload: respData });
        }
      } catch (error) {
        console.error(error);
        dispatch({ type: "FETCH_ERROR", payload: error });
      }
    },
    []
  );

  useEffect(() => {
    fetchOnMount && fetchPlexWebhooks({ pageIndex, pageSize, orderBy, sortBy });
  }, [fetchOnMount, fetchPlexWebhooks, orderBy, pageIndex, pageSize, sortBy]);

  return {
    ...state,
    fetchPlexWebhooks,
  };
}
