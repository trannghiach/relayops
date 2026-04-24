"use client";

import { api } from "@/lib/api";
import type { DeadLetter } from "@/lib/types";
import { Button, message, Table, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import Link from "next/link";
import { useEffect, useState } from "react";

export default function DeadLettersPage() {
  const [items, setItems] = useState<DeadLetter[]>([]);
  const [loading, setLoading] = useState(true);

  async function load() {
    setLoading(true);
    try {
      const res = await api.getDeadLetters();
      setItems(res.data);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
  }, []);

  async function handleReplay(jobID: string) {
    await api.replayJob(jobID);
    message.success("Job replayed");
    await load();
  }

  const columns: ColumnsType<DeadLetter> = [
    {
      title: "Job ID",
      dataIndex: "job_id",
      render: (v) => <Link href={`/jobs/${v}`}>{v}</Link>,
    },
    {
      title: "Reason",
      dataIndex: "reason",
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      render: (v) => new Date(v).toLocaleString(),
    },
    {
      title: "Action",
      render: (_, row) => (
        <Button danger size="small" onClick={() => handleReplay(row.job_id)}>
          Replay
        </Button>
      ),
    },
  ];

  return (
    <>
      <Typography.Title level={2}>Dead Letters</Typography.Title>

      <Table
        rowKey="id"
        loading={loading}
        columns={columns}
        dataSource={items}
        pagination={{ pageSize: 10 }}
      />
    </>
  );
}