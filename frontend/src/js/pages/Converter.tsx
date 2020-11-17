import * as React from "react";
import {useRef, useState} from "react";
import {Button} from "@app/components";

export const Converter = () => {
    const [code, setCode] = useState("");
    const textAreaRef = useRef<any>();

    const onConvert = () => {
        (window as any).backend.convert(code).then((code: string) => {
            setCode(code);
            textAreaRef.current.select();
            document.execCommand('copy');
        });
    }

    return (
        <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="border-b border-gray-200 px-4 py-5 sm:px-6">
                <div className="flex justify-end">
                    <Button
                        onClick={onConvert}
                    >
                        Convert
                    </Button>
                </div>
            </div>

            <div className="bg-gray-50 px-4 py-5 sm:p-6">
                <div className="space-y-1">
                    <div className="rounded-md shadow-sm">
                        <textarea
                            ref={textAreaRef}
                            rows={15}
                            value={code}
                            onChange={e => setCode(e.target.value)}
                            className="form-textarea mt-1 block w-full transition duration-150 ease-in-out sm:text-sm sm:leading-5">
                        </textarea>
                    </div>
                </div>
            </div>
        </div>
    )
}
