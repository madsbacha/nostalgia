"use client"

import { zodResolver } from "@hookform/resolvers/zod"
import {Noop, useForm} from "react-hook-form"
import { z } from "zod"

import { Button } from "@/components/ui/button"
import {
    Form,
    FormControl, FormDescription,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import {Textarea} from "@/components/ui/textarea";
import {uploadMediaFile} from "@/lib/api";
import {redirect} from "next/navigation";

const parseTags: (tags: string) => string[] = tags => tags.split(",").map(tag => tag.trim())

const FormSchema = z.object({
    title: z.string().min(2, {
        message: "Title must be at least 2 characters.",
    }),
    description: z.string().max(1000, {
        message: "Description must be at most 1000 characters.",
    }),
    file: z.instanceof(File).refine(
        (file) =>
            [
                "video/mp4",
                "video/webm",
                "video/ogg",
                "video/x-msvideo",     // AVI
                "video/quicktime",     // MOV
                "video/x-matroska"     // MKV
            ].includes(file.type),
        { message: "Invalid file type" }
    ),
    tags: z.string()
        .refine(tags => {
            const regex = /^[a-zA-Z0-9_]+$/;
            return parseTags(tags).every(tag => regex.test(tag));
        }, { message: "Tags should only contain alphanumeric characters and underscores" })
        .refine(tags => {
            const tagsArray = parseTags(tags);
            return new Set(tagsArray).size === tagsArray.length;
        }, { message: "Tags should be unique" })
        .refine(
            tags => parseTags(tags).every(tag => tag.length <= 20),
            { message: "Tags should be at most 20 characters long" })
        .refine(
            tags => parseTags(tags).length <= 10,
            { message: "You can only add up to 10 tags" })
})

interface InputFileProps {
    value: File
    onChange: (file: File) => void
    accept: string
    disabled?: boolean | undefined;
    name: string
    ref: React.Ref<HTMLInputElement>
    onBlur: Noop
}

function InputFile({ accept, disabled, name, ref, onBlur, onChange }: InputFileProps) {
    const onChangeHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0]
        if (file) {
            console.log(file)
            if (onChange) {
                onChange(file)
            }
        }
    }

    return (
        <Input type="file"
                accept={accept}
                disabled={disabled}
                name={name}
                ref={ref}
                onBlur={onBlur}
                onChange={onChangeHandler}
        />
    )
}

function uploadForm(data: z.infer<typeof FormSchema>) {
    uploadMediaFile(data.file, data.title, data.description, data.tags)
        .then(r => redirect(`/media/${r.media_id}`))
}

export function UploadMediaForm() {
    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues: {
            title: "",
            description: "",
            tags: "",
        },
    })

    function onSubmit(data: z.infer<typeof FormSchema>) {
        // toast({
        //     title: "You submitted the following values:",
        //     description: (
        //         <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
        //   <code className="text-white">{JSON.stringify(data, null, 2)}</code>
        // </pre>
        //     ),
        // })
        console.log(data)
        uploadForm(data)
    }

    return (
        <>
            <h1 className="text-2xl py-6">Upload</h1>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="w-full space-y-6">
                    <FormField
                        control={form.control}
                        name="title"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Title</FormLabel>
                                <FormControl>
                                    <Input {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="description"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Description</FormLabel>
                                <FormControl>
                                    <Textarea {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="file"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Description</FormLabel>
                                <FormControl>
                                    <InputFile {...field} accept="video/*" />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <FormField
                        control={form.control}
                        name="tags"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Tags</FormLabel>
                                <FormControl>
                                    <Input {...field} />
                                </FormControl>
                                <FormDescription>A comma-separated list of tags</FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />
                    <Button type="submit" disabled={form.formState.isSubmitting}>Submit</Button>
                </form>
            </Form>
        </>
    )
}
