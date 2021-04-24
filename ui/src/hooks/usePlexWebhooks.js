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
        ...state,
        loading: false,
        plexWebhooks: action.payload,
        total: action.payload.total ?? -1,
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

function buildURL({
  pageIndex,
  pageSize,
  orderBy = initialState.pagination.orderBy,
  sortBy = initialState.pagination.sortBy,
}) {
  const baseURL = "/plex";
  return pageIndex >= 0 && pageSize !== -1
    ? `${baseURL}?limit=${pageSize}&offset=${
        pageIndex * pageSize
      }&orderBy=${orderBy}&sortBy=${sortBy}`
    : baseURL;
}

export function usePlexWebhooks({
  fetchOnMount = false,
  pageIndex = initialState.pagination.pageIndex,
  pageSize = initialState.pagination.pageSize,
  orderBy = initialState.pagination.orderBy,
  sortBy = initialState.pagination.sortBy,
} = {}) {
  const [state, dispatch] = useReducer(
    reducer,
    {
      loading: fetchOnMount,
      pagination: {
        pageIndex,
        pageSize,
        orderBy,
        sortBy,
      },
    },
    (args) => ({
      ...initialState,
      loading: args.loading,
      pagination: {
        ...initialState.pagination,
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

  const setPagination = useCallback((paginationOptions) => {
    dispatch({ type: "UPDATE_PAGINATION", payload: paginationOptions });
  }, []);

  useEffect(() => {
    fetchOnMount && fetchPlexWebhooks({ pageIndex, pageSize, orderBy, sortBy });
  }, [fetchOnMount, fetchPlexWebhooks, orderBy, pageIndex, pageSize, sortBy]);

  return {
    ...state,
    fetchPlexWebhooks,
    setPagination,
  };
}
