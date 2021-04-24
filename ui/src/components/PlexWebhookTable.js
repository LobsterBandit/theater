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
      pageCount: total === -1 ? total : total / pagination.pageSize,
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
        count={
          total === -1 && plexWebhooks.length === pageSize
            ? total
            : pageIndex * pageSize + plexWebhooks.length
        }
        labelDisplayedRows={({ from, to, count, page }) => {
          const pageText = `Page: ${page + 1} of ${
            count !== -1 ? page + 1 : `more than ${page + 1}`
          }`;

          const rowText = `Rows: ${from}-${count !== -1 ? count : to} of ${
            count !== -1 ? count : `more than ${to}`
          }`;

          return `${pageText} | ${rowText}`;
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
