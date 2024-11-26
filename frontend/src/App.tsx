import React, { useState, useEffect } from "react";
import Modal from "./Modal";
import { CheckForUpdate, Click, Reset, GetTimings, ToggleRounding, GetRoundState, CycleDivisionMode, GetDivisionMode } from "../wailsjs/go/main/App";

interface Timing {
    Full: string
    Half: string
    Quarter: string
    Eighth: string
    Sixteenth: string
    ThirtySecond: string
    SixtyFourth: string
    OneTwentyEighth: string
}

const App: React.FC = () => {
    const [updateInfo, setUpdateInfo] = useState<{ available: boolean; message: string; url?: string } | null>(null);
    const [showModal, setShowModal] = useState<boolean>(false);

    const [bpm, setBpm] = useState<number>(0);
    const [timing, setTiming] = useState<Timing | null>(null);
    const [roundOutputs, setRoundOutputs] = useState<boolean>(false);
    const [divisionMode, setDivisionMode] = useState<number>(0); // 0: NoDivision, 1: DivideBy100, 2: DivideBy1000

    // Fetch update info on mount
    useEffect(() => {
        const fetchUpdateInfo = async () => {
            try {
                const info = await CheckForUpdate();
                setUpdateInfo(info);
                if (info.available) {
                    setShowModal(true);
                }
            } catch (err) {
                console.error("Error checking for updates:", err);
            }
        };

        fetchUpdateInfo();
    }, []);

    // Initialize app state on mount
    useEffect(() => {
        const initializeAppState = async () => {
            try {
                const timings: Timing = await GetTimings();
                const isRounded: boolean = await GetRoundState();
                const mode: number = await GetDivisionMode();

                setTiming(timings);
                setRoundOutputs(isRounded);
                setDivisionMode(mode);
            } catch (err) {
                console.error("Error initializing app state:", err);
            }
        };

        initializeAppState();
    }, []);

    const handleTap = async () => {
        const newBpm: number = await Click();
        setBpm(parseFloat(newBpm.toFixed(2)));

        const timings: Timing = await GetTimings();
        setTiming(timings);
    };

    const handleReset = async () => {
        await Reset();
        setBpm(0);
        setTiming(null);
    };

    const handleToggleRounding = async () => {
        await ToggleRounding();
        const timings: Timing = await GetTimings();
        const isRounded: boolean = await GetRoundState();
        setTiming(timings);
        setRoundOutputs(isRounded);
    };

    const handleCycleDivisionMode = async () => {
        await CycleDivisionMode();
        const timings: Timing = await GetTimings();
        const mode: number = await GetDivisionMode();
        setTiming(timings);
        setDivisionMode(mode);
    };

    useEffect(() => {
        const handleKeyPress = async (event: KeyboardEvent) => {
            if (event.key === "r" || event.key === "R") {
                handleReset();
            } else if (event.key === "F1") {
                event.preventDefault();
                await handleToggleRounding();
            } else if (event.key === "F2") {
                event.preventDefault();
                await handleCycleDivisionMode();
            }
        };

        window.addEventListener("keydown", handleKeyPress);
        return () => {
            window.removeEventListener("keydown", handleKeyPress);
        };
    }, []);

    return (
        <div className="container">
            <Modal
                show={showModal}
                title="Update Available"
                message={updateInfo?.message || ""}
                url={updateInfo?.url}
                onClose={() => setShowModal(false)}
            />

            <div className="circle" onClick={handleTap}>
                <span className="bpm-text">{bpm || "Tap"}</span>
            </div>
            <p className="instructions">
                Press <strong>R</strong> to reset, <strong>F1</strong> to {roundOutputs ? "disable" : "enable"} rounding, <strong>F2</strong> to cycle division mode
                <span className="time-header-asterisk">
                    ({divisionMode === 0
                        ? "none"
                        : divisionMode === 1
                            ? "÷ 100"
                            : "÷ 1000"})
                </span>
                .
            </p>
            <div className={`table-container ${timing ? "visible" : "hidden"}`}>
                {timing && (
                    <table className="table">
                        <thead>
                            <tr>
                                <th>Note</th>
                                <th>Time</th>
                            </tr>
                        </thead>
                        <tbody>
                            {[
                                ["1", timing.Full],
                                ["1/2", timing.Half],
                                ["1/4", timing.Quarter],
                                ["1/8", timing.Eighth],
                                ["1/16", timing.Sixteenth],
                                ["1/32", timing.ThirtySecond],
                                ["1/64", timing.SixtyFourth],
                                ["1/128", timing.OneTwentyEighth],
                            ].map(([note, value], idx) => (
                                <tr key={idx}>
                                    <td>{note}</td>
                                    <td>{value}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                )}
            </div>
        </div>
    );
};


export default App
