import {
  Table,
  TableBody,
  TableContainer,
  TableCell,
  TableHead,
  TablePagination,
  TableRow,
} from "@material-ui/core";
import TablePaginationActions from "@material-ui/core/TablePagination/TablePaginationActions";
import { useEffect } from "react";
import { usePagination, useTable } from "react-table";
import { PlexWebhookToolbar } from "./PlexWebhookToolbar";
import { usePlexWebhooks } from "../hooks/usePlexWebhooks";

const columns = [
  { Header: "ID", accessor: "id" },
  { Header: "Date", accessor: "date" },
  { Header: "Event", accessor: "payload.event" },
  { Header: "User", accessor: "payload.Account.title" },
  { Header: "Player", accessor: "payload.Player.title" },
  { Header: "Type", accessor: "payload.Metadata.type" },
  { Header: "Title", accessor: "payload.Metadata.title" },
  {
    Header: "Payload",
    accessor: "payload",
    Cell: ({ value }) => (
      <pre
        style={{
          whiteSpace: "pre-wrap",
          wordWrap: "break-word",
          background: "lightgray",
        }}
      >
        {JSON.stringify(value, null, 2)}
      </pre>
    ),
  },
];

export function PlexWebhookTable() {
  const {
    fetchPlexWebhooks,
    pagination,
    plexWebhooks,
    setPagination,
    total,
  } = usePlexWebhooks();

  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    page,
    prepareRow,

    // pagination
    pageCount,
    gotoPage,
    setPageSize,
    state: { pageIndex, pageSize },
  } = useTable(
    {
      columns,
      data: plexWebhooks,
      initialState: {
        hiddenColumns: ["payload"],
        pageIndex: pagination.pageIndex,
        pageSize: pagination.pageSize,
      },
      manualPagination: true,
      pageCount: total === -1 ? total : total / 10,
    },
    usePagination
  );

  useEffect(() => {
    fetchPlexWebhooks({ pageIndex, pageSize });
  }, [fetchPlexWebhooks, pageIndex, pageSize]);

  return (
    <>
      <PlexWebhookToolbar
        onRefreshClick={() => fetchPlexWebhooks(pagination)}
      />
      <TableContainer sx={{ maxHeight: "75vh" }}>
        <Table component="div" size="small" stickyHeader {...getTableProps()}>
          <TableHead component="div">
            {headerGroups.map((headerGroup) => (
              <TableRow component="div" {...headerGroup.getHeaderGroupProps()}>
                {headerGroup.headers.map((column) => (
                  <TableCell component="div" {...column.getHeaderProps()}>
                    {column.render("Header")}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableHead>
          <TableBody component="div" {...getTableBodyProps()}>
            {page.map((row) => {
              prepareRow(row);
              return (
                <TableRow component="div" hover {...row.getRowProps()}>
                  {row.cells.map((cell) => {
                    return (
                      <TableCell component="div" {...cell.getCellProps()}>
                        {cell.render("Cell")}
                      </TableCell>
                    );
                  })}
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        ActionsComponent={({ ...props }) => (
          <TablePaginationActions
            {...props}
            showFirstButton={true}
            showLastButton={true}
          />
        )}
        component="div"
        count={total === -1 ? total : plexWebhooks.length}
        labelDisplayedRows={({ from, to, count, page }) => {
          console.log({ page, pageCount, count, from, to });
          return `Page: ${page + 1} of ${
            pageCount === 0
              ? 1
              : count !== -1
              ? pageCount
              : `more than ${page + 1}`
          } | Rows: ${from}-${to} of ${
            count !== -1 ? count : `more than ${to}`
          }`;
        }}
        onRowsPerPageChange={(e) => {
          const pageSize = Number(e.target.value);
          setPagination({ pageSize });
          setPageSize(pageSize);
        }}
        onPageChange={(e, page) => {
          setPagination({ pageIndex: page });
          gotoPage(page);
        }}
        page={pageIndex}
        rowsPerPage={pageSize}
        rowsPerPageOptions={[10, 25, 50, { label: "All", value: -1 }]}
      />
    </>
  );
}
