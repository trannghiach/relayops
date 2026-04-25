"use client";

import { api } from "@/lib/api";
import type { EventItem } from "@/lib/types";
import { Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useEffect, useState } from "react";

export default function EventsPage() {
  const [events, setEvents] = useState<EventItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getEvents()
      .then((res) => setEvents(res.data))
      .finally(() => setLoading(false));
  }, []);

  const columns: ColumnsType<EventItem> = [
    {
      title: "Event Type",
      dataIndex: "event_type",
    },
    {
      title: "Source",
      dataIndex: "source",
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (status) => (
        <Tag color={status === "planned" ? "green" : "red"}>{status}</Tag>
      ),
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      render: (v) => new Date(v).toLocaleString(),
    },
    {
      title: "ID",
      dataIndex: "id",
      render: (v) => <Typography.Text copyable>{v}</Typography.Text>,
    },
  ];

  return (
    <>
      <Typography.Title level={2}>Events</Typography.Title>

      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={events}
        pagination={{ pageSize: 10 }}
        scroll={{ x: "max-content" }}
      />
    </>
  );
}
