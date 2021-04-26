import {
  Table,
  TableBody,
  TableContainer,
  TableCell,
  TableHead,
  TablePagination,
  TableRow,
  IconButton,
  Tooltip,
} from "@material-ui/core";
import TablePaginationActions from "@material-ui/core/TablePagination/TablePaginationActions";
import Replay from "@material-ui/icons/Replay";
import { useEffect, useState } from "react";
import { usePagination, useTable } from "react-table";
import { PlexWebhookToolbar } from "./PlexWebhookToolbar";
import { WebhookPayloadDialog } from "./WebhookPayloadDialog";
import { usePlexWebhooks } from "../hooks/usePlexWebhooks";
import { replayPlexWebhook } from "../api";

const columns = [
  { Header: "ID", accessor: "id" },
  { Header: "Date", accessor: "date" },
  { Header: "Event", accessor: "payload.event" },
  { Header: "User", accessor: "payload.Account.title" },
  { Header: "Player", accessor: "payload.Player.title" },
  { Header: "Type", accessor: "payload.Metadata.type" },
  { Header: "Title", accessor: "payload.Metadata.title" },
  {
    Header: "Replay",
    id: "replay",
    Cell: ({ row: { original } }) => {
      return (
        <Tooltip placement="left" title="Replay event">
          <IconButton size="small">
            <Replay />
          </IconButton>
        </Tooltip>
      );
    },
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

  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectedPayload, setSelectedPayload] = useState({});

  const openDialog = (data) => {
    setSelectedPayload(data);
    setDialogOpen(true);
  };

  const handleDialogClose = () => {
    setSelectedPayload({});
    setDialogOpen(false);
  };

  const handleReplay = (e, payload) => {
    e.stopPropagation();
    replayPlexWebhook(payload)
      .then((r) => console.log(r))
      .catch((e) => console.error(e));
  };

  useEffect(() => {
    fetchPlexWebhooks({ pageIndex, pageSize });
  }, [fetchPlexWebhooks, pageIndex, pageSize]);

  return (
    <>
      <WebhookPayloadDialog
        handleClose={handleDialogClose}
        handleReplay={handleReplay}
        open={dialogOpen}
        value={selectedPayload}
      />
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
                <TableRow
                  component="div"
                  hover
                  onClick={() => openDialog(row.original)}
                  {...row.getRowProps()}
                >
                  {row.cells.map((cell) => {
                    return (
                      <TableCell
                        component="div"
                        {...cell.getCellProps({
                          ...(cell.column.id === "replay" && {
                            style: {
                              padding: 0,
                              textAlign: "center",
                            },
                            onClick: (e) => {
                              e.stopPropagation();
                              handleReplay(e, row.original.payload);
                            },
                          }),
                        })}
                      >
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
