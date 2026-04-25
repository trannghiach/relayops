"use client";

import {
    DashboardOutlined,
    DatabaseOutlined,
    MenuOutlined,
    WarningOutlined,
    ProfileOutlined,
} from "@ant-design/icons";
import { Button, Drawer, Grid, Layout, Menu, Space, Typography } from "antd";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useState, type ReactNode } from "react";

const { Header, Sider, Content } = Layout;

export function AppShell({ children }: { children: ReactNode }) {
    const pathname = usePathname();
    const screens = Grid.useBreakpoint();
    const isMobile = !screens.md;
    const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

    const menuItems = [
        {
            key: "/",
            icon: <DashboardOutlined />,
            label: <Link href="/">Overview</Link>,
        },
        {
            key: "/events",
            icon: <ProfileOutlined />,
            label: <Link href="/events">Events</Link>,
        },
        {
            key: "/jobs",
            icon: <DatabaseOutlined />,
            label: <Link href="/jobs">Jobs</Link>,
        },
        {
            key: "/dead-letters",
            icon: <WarningOutlined />,
            label: <Link href="/dead-letters">Dead Letters</Link>,
        },
    ];

    return (
        <Layout style={{ minHeight: "100vh" }}>
            {!isMobile && (
                <Sider width={240}>
                    <div style={{ padding: 20 }}>
                        <Typography.Title level={4} style={{ color: "white", margin: 0 }}>
                            RelayOps
                        </Typography.Title>
                    </div>

                    <Menu theme="dark" mode="inline" selectedKeys={[pathname]} items={menuItems} />
                </Sider>
            )}

            <Layout>
                <Header style={{ background: "white", paddingInline: isMobile ? 12 : 24 }}>
                    <Space size="middle">
                        {isMobile && (
                            <Button
                                aria-label="Open menu"
                                icon={<MenuOutlined />}
                                onClick={() => setMobileMenuOpen(true)}
                            />
                        )}
                        <Typography.Text strong>
                            Event-driven notification & workflow platform
                        </Typography.Text>
                    </Space>
                </Header>

                <Content style={{ padding: isMobile ? 12 : 24 }}>{children}</Content>
            </Layout>

            <Drawer
                title="RelayOps"
                placement="left"
                open={mobileMenuOpen}
                onClose={() => setMobileMenuOpen(false)}
                bodyStyle={{ padding: 0 }}
            >
                <Menu
                    mode="inline"
                    selectedKeys={[pathname]}
                    items={menuItems}
                    onClick={() => setMobileMenuOpen(false)}
                />
            </Drawer>
        </Layout>
    );
}
