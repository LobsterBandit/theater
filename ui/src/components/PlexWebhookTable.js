import {
  Table,
  TableBody,
  TableContainer,
  TableCell,
  TableHead,
  TableRow,
  TablePagination,
} from "@material-ui/core";
import TablePaginationActions from "@material-ui/core/TablePagination/TablePaginationActions";
import { useEffect } from "react";
import { usePagination, useTable } from "react-table";

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

export function PlexWebhookTable({ data, fetchData, loading, totalCount }) {
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
      data,
      initialState: { hiddenColumns: ["payload"], pageIndex: 0, pageSize: 10 },
      manualPagination: true,
      pageCount: totalCount === -1 ? totalCount : totalCount / 10,
    },
    usePagination
  );

  useEffect(() => {
    fetchData({ pageIndex, pageSize });
  }, [fetchData, pageIndex, pageSize]);

  return (
    <>
      <TableContainer sx={{ maxHeight: "75vh" }}>
        <Table size="small" stickyHeader {...getTableProps()}>
          <TableHead>
            {headerGroups.map((headerGroup) => (
              <TableRow {...headerGroup.getHeaderGroupProps()}>
                {headerGroup.headers.map((column) => (
                  <TableCell {...column.getHeaderProps()}>
                    {column.render("Header")}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableHead>
          <TableBody {...getTableBodyProps()}>
            {page.map((row) => {
              prepareRow(row);
              return (
                <TableRow hover {...row.getRowProps()}>
                  {row.cells.map((cell) => {
                    return (
                      <TableCell {...cell.getCellProps()}>
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
        count={totalCount === -1 ? totalCount : data.length}
        labelDisplayedRows={({ from, to, count, page }) =>
          `Page: ${page + 1} of ${
            pageCount === 0
              ? 1
              : count !== -1
              ? pageCount
              : `more than ${page + 1}`
          } | Rows: ${from}-${to} of ${
            count !== -1 ? count : `more than ${to}`
          }`
        }
        onRowsPerPageChange={(e) => setPageSize(Number(e.target.value))}
        onPageChange={(e, page) => gotoPage(page)}
        page={pageIndex}
        rowsPerPage={pageSize}
        rowsPerPageOptions={[10, 25, 50, { label: "All", value: -1 }]}
      />
    </>
  );
}
