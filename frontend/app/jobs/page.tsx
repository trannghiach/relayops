"use client";

import { api } from "@/lib/api";
import type { JobItem } from "@/lib/types";
import { Button, Table, Tag, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import Link from "next/link";
import { useEffect, useState } from "react";

function statusColor(status: string) {
  switch (status) {
    case "succeeded":
      return "green";
    case "pending":
      return "blue";
    case "processing":
      return "orange";
    case "dead_lettered":
      return "red";
    case "failed":
      return "volcano";
    default:
      return "default";
  }
}

export default function JobsPage() {
  const [jobs, setJobs] = useState<JobItem[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getJobs()
      .then((res) => setJobs(res.data))
      .finally(() => setLoading(false));
  }, []);

  const columns: ColumnsType<JobItem> = [
    {
      title: "Job Type",
      dataIndex: "job_type",
    },
    {
      title: "Channel",
      dataIndex: "channel",
    },
    {
      title: "Status",
      dataIndex: "status",
      render: (status) => <Tag color={statusColor(status)}>{status}</Tag>,
    },
    {
      title: "Attempts",
      render: (_, row) => `${row.attempts}/${row.max_attempts}`,
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      render: (v) => new Date(v).toLocaleString(),
    },
    {
      title: "Action",
      render: (_, row) => (
        <Link href={`/jobs/${row.id}`}>
          <Button size="small">View</Button>
        </Link>
      ),
    },
  ];

  return (
    <>
      <Typography.Title level={2}>Jobs</Typography.Title>

      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={jobs}
        pagination={{ pageSize: 10 }}
      />
    </>
  );
}